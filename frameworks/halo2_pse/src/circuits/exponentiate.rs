use super::is_zero::{IsZeroChip, IsZeroConfig};
use super::utils::f_values::{get_one, get_zero};
use std::marker::PhantomData;
use halo2_proofs::{
    arithmetic::Field,
    circuit::{AssignedCell, Layouter, SimpleFloorPlanner, Value},
    plonk::{Advice, Circuit, Column, ConstraintSystem, Error, Instance, Selector, Expression},
    poly::{Rotation}, halo2curves::bn256::Fr, 
};
use serde::{Deserialize, Serialize};
use serde_json;

#[derive(Debug, Clone)]
pub struct ExponentiationConfig<F: Field>  {
    col_x: Column<Advice>,
    col_e: Column<Advice>,
    col_y: Column<Advice>,
    e_equals_zero: IsZeroConfig<F>,
    selector: Selector,
    instance: Column<Instance>,
}

#[derive(Debug, Clone)]
pub struct ExponentiationChip<F: Field> {
    config: ExponentiationConfig<F>,
    _marker: PhantomData<F>,
}

impl<F: Field> ExponentiationChip<F> {
    pub fn construct(config: ExponentiationConfig<F>) -> Self {
        Self {
            config,
            _marker: PhantomData,
        }
    }

    pub fn configure(
        meta: &mut ConstraintSystem<F>,
    ) -> ExponentiationConfig<F> {
        let selector = meta.selector();
        let col_x = meta.advice_column();
        let col_e = meta.advice_column();
        let col_y = meta.advice_column();
        let is_zero_advice_column = meta.advice_column();
        let instance = meta.instance_column();

        meta.enable_equality(col_x);
        meta.enable_equality(col_e);
        meta.enable_equality(col_y);
        meta.enable_equality(instance);

        let e_equals_zero = IsZeroChip::configure(
           meta,
           |meta| meta.query_selector(selector), 
           |meta| meta.query_advice(col_e, Rotation::cur()),
           is_zero_advice_column, // this is the advice column that stores value_inv
       );

        meta.create_gate("exponentiate", |meta| {
            // y == x**e, x = 2 e = 3 y = 8

            //     x    |    e      |    y    |    s
            //
            //     2    |    0      |    1
            //     2    |    1      |    2
            //     2    |    2      |    4
            //
            // e + 1 - e_next = 0
            //
            // is_zero(e)(1 - y) + (1-is_zero(e))(x * y - y_next) = 0
            let s = meta.query_selector(selector);
            let x = meta.query_advice(col_x, Rotation::cur());
            let e = meta.query_advice(col_e, Rotation::cur());
            let e_next = meta.query_advice(col_e, Rotation::next());
            let y = meta.query_advice(col_y, Rotation::cur());
            let y_next = meta.query_advice(col_y, Rotation::next());

            let one = Expression::Constant(get_one());
            let iszero = e_equals_zero.expr();
            vec![
                //s.clone() * (e + one.clone() - e_next),
                s.clone() * (iszero.clone() * (one.clone() - y.clone()) + (one.clone() - iszero) * (x * y.clone() - y_next))
            ]
        });

        ExponentiationConfig {
            col_x,
            col_e,
            col_y,
            e_equals_zero,
            selector,
            instance,
        }
    }

    pub fn assign(
        &self,
        mut layouter: impl Layouter<F>,
        nrows: usize,
    ) -> Result<AssignedCell<F, F>, Error> {
        let is_zero_chip = IsZeroChip::construct(self.config.e_equals_zero.clone());
        layouter.assign_region(
            || "entire exponentiation table",
            |mut region| {
                self.config.selector.enable(&mut region, 0)?;

                let mut x_cell = region.assign_advice_from_instance(
                    || "x",
                    self.config.instance,
                    0,
                    self.config.col_x,
                    0,
                )?;
                let mut e_cell = region.assign_advice(
                    || "e",
                    self.config.col_e,
                    0,
                    || -> Value<F> {Value::known(get_zero())},
                )?;
                let value_e: Value<F> = e_cell.value().clone().map(|f_ref| f_ref.clone()); 
                is_zero_chip.assign(&mut region, 0, value_e)?;
                let mut y_cell = region.assign_advice(
                    || "y",
                    self.config.col_y,
                    0,
                    || -> Value<F> {Value::known(get_one())},
                )?;

                let mut row = 1;
                while row <= nrows {
                    if row < nrows - 1 {
                        self.config.selector.enable(&mut region, row)?;
                    }
                    x_cell = x_cell.copy_advice(
                        || "x", 
                        &mut region, 
                        self.config.col_x, 
                        row
                    )?;
                    e_cell = region.assign_advice(
                        || "e",
                        self.config.col_e,
                        row,
                        || e_cell.value().copied() + Value::known(get_one::<F>())
                    )?;
                    y_cell = region.assign_advice(
                        || "y",
                        self.config.col_y,
                        row,
                        || y_cell.value().copied() * x_cell.value().copied()
                    )?;
                    let value_e: Value<F> = e_cell.value().clone().map(|f_ref| f_ref.clone()); 
                    is_zero_chip.assign(&mut region, row, value_e)?;
                    row += 1;
                }

                Ok(y_cell)
            },
        )
    }

    pub fn expose_public(
        &self,
        mut layouter: impl Layouter<F>,
        cell: AssignedCell<F, F>,
        row: usize,
    ) -> Result<(), Error> {
        layouter.constrain_instance(cell.cell(), self.config.instance, row)
    }
}

#[derive(Default, Clone)]
pub struct ExponentiationCircuit {
    pub row: usize,
}

impl<F: Field> Circuit<F> for ExponentiationCircuit {
    type Config = ExponentiationConfig<F>;
    type FloorPlanner = SimpleFloorPlanner;

    fn without_witnesses(&self) -> Self {
        Self::default()
    }

    fn configure(meta: &mut ConstraintSystem<F>) -> Self::Config {
        ExponentiationChip::configure(meta)
    }

    fn synthesize(
        &self,
        config: Self::Config,
        mut layouter: impl Layouter<F>,
    ) -> Result<(), Error> {
        let chip = ExponentiationChip::construct(config);

        let out_cell = chip.assign(layouter.namespace(|| "entire table"), self.row)?;

        chip.expose_public(layouter.namespace(|| "out"), out_cell, 2)?;

        Ok(())
    }
}

#[derive(Debug, Deserialize, Serialize)]
pub struct ExponentiationData {
    X: String,
    E: String,
    Y: String
}


pub fn get_exponentiation_data (
    input_file_str: String
) -> (u32, usize, Fr, Fr, Fr){
    let data: ExponentiationData = serde_json::from_str(&input_file_str).expect("Cannot read json string");
    let x_value: usize = data.X.parse().unwrap();
    let e_value: usize = data.E.parse().unwrap();
    let y_value: usize = data.Y.parse().unwrap();
    let x = Fr::from(x_value as u64); 
    let e = Fr::from(e_value as u64); 
    let y = Fr::from(y_value as u64); 
    let k = match e_value {
        0..=31 => 5,
        32..=63 => 6,
        64..=127 => 7,
        128..=255 => 8,
        256..=511 => 9,
        512..=1023 => 10,
        1024..=2047 => 11,
        2048..=4095 => 12,
        4096..=8191 => 13,
        8192..=16383 => 14,
        16384..=32767 => 15,
        32768..=65535 => 16,
        65536..=131071 => 17,
        131072..=262143 => 18,
        262144..=524287 => 19,
        524288..=1048575 => 20,
        1048576..=2097151 => 21,
        2097152..=4194303 => 22,
        4194304..=8388607 => 23,
        8388608..=16777215 => 24,
        _ => 25,
    };
    return (k, e_value, x, e, y);
}


#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_exponentiation_circuit() {
        use halo2_proofs::dev::MockProver;
        let k = 5;
        let e_value: usize = 12;
        let circuit = ExponentiationCircuit {
            row: e_value,
        };
        let x = Fr::from(2); 
        let e = Fr::from(e_value as u64); 
        let y = Fr::from(4096); 

        //let public_input = vec![x, e, y];
        let public_input = vec![x, e, y];
        let prover = MockProver::run(k, &circuit, vec![public_input.clone()]).unwrap();
        prover.assert_satisfied();
    }
}
